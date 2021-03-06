import 'package:angular/angular.dart';
import 'package:angular_forms/angular_forms.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/gravatar_component.dart';
import 'package:gitpods/mailto_pipe.dart';
import 'package:gitpods/repository.dart';
import 'package:gitpods/user.dart';
import 'package:gitpods/user_service.dart';

@Component(
  selector: 'gitpods-user-profile',
  templateUrl: 'user_profile_component.html',
  styleUrls: const ['user_profile_component.css'],
  providers: const [UserService],
  directives: const [
    COMMON_DIRECTIVES, ROUTER_DIRECTIVES, formDirectives, Gravatar
  ],
  pipes: const [DatePipe, MailtoPipe, FilteredReposPipe],
)
class UserProfileComponent implements OnInit {
  final RouteParams _routeParams;
  final UserService _userService;

  UserProfileComponent(this._routeParams, this._userService);

  User user;
  String repoQuery = '';

  @override
  void ngOnInit() {
    String username = this._routeParams.get('username');
    this._userService.profile(username)
        .then((user) => this.user = user);
  }
}

@Pipe('filteredRepos')
class FilteredReposPipe extends PipeTransform {
  List<Repository> transform(List<Repository> repos, String query) {
    repos = repos.where((repo) => repo.name.contains(query)).toList();
    repos.sort((a, b) => a.name.compareTo(b.name));
    return repos;
  }
}
